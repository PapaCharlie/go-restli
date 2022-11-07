package io.papacharlie.gorestli;

import com.google.common.base.Preconditions;
import com.google.common.collect.ImmutableSet;
import com.linkedin.data.schema.ArrayDataSchema;
import com.linkedin.data.schema.DataSchema;
import com.linkedin.data.schema.DataSchemaLocation;
import com.linkedin.data.schema.EnumDataSchema;
import com.linkedin.data.schema.FixedDataSchema;
import com.linkedin.data.schema.MapDataSchema;
import com.linkedin.data.schema.NamedDataSchema;
import com.linkedin.data.schema.RecordDataSchema;
import com.linkedin.data.schema.TyperefDataSchema;
import com.linkedin.data.schema.UnionDataSchema;
import com.linkedin.pegasus.generator.DataSchemaParser;
import com.linkedin.restli.restspec.CollectionSchema;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestSpecCodec;
import com.linkedin.restli.tools.snapshot.gen.SnapshotGenerator;
import io.papacharlie.gorestli.json.ComplexKey;
import io.papacharlie.gorestli.json.DataType;
import io.papacharlie.gorestli.json.Enum;
import io.papacharlie.gorestli.json.Fixed;
import io.papacharlie.gorestli.json.PathKey;
import io.papacharlie.gorestli.json.Record;
import io.papacharlie.gorestli.json.Record.Field;
import io.papacharlie.gorestli.json.RestliType;
import io.papacharlie.gorestli.json.StandaloneUnion;
import io.papacharlie.gorestli.json.Typeref;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import javax.annotation.Nullable;

import static com.linkedin.data.schema.DataSchema.Type.NULL;
import static com.linkedin.data.schema.DataSchema.Type.RECORD;
import static com.linkedin.data.schema.DataSchema.Type.TYPEREF;
import static com.linkedin.data.schema.DataSchema.Type.UNION;
import static io.papacharlie.gorestli.json.RestliType.GoPrimitive;
import static io.papacharlie.gorestli.json.RestliType.Identifier;
import static io.papacharlie.gorestli.json.RestliType.JAVA_TO_GO_PRIMTIIVE_TYPE;
import static io.papacharlie.gorestli.json.RestliType.RAW_RECORD;
import static io.papacharlie.gorestli.json.RestliType.UnionMember;


public class TypeParser {
  // These types are manually defined in go-restli and should never be re-generated
  private static final Set<Identifier> RESERVED_TYPES = ImmutableSet.of(
      MethodParser.PAGING_CONTEXT,
      new Identifier("com.linkedin.restli.common", "EmptyRecord"),
      new Identifier("restlidata", "RawRecord")
  );

  private final DataSchemaParser _parser;
  private final Set<String> _rawRecords;
  private final Map<Identifier, DataType> _dataTypes = new HashMap<>();
  private final Set<String> _anonymousTypeRefs = new HashSet<>();

  public TypeParser(DataSchemaParser parser, Set<String> rawRecords) {
    _parser = parser;
    _rawRecords = rawRecords;
  }

  public void extractDataTypes(ResourceSchema resourceSchema) {
    SnapshotGenerator generator = new SnapshotGenerator(resourceSchema, _parser.getSchemaResolver());
    for (NamedDataSchema schema : generator.generateModelList()) {
      addNamedDataSchema(schema.getFullName());
    }
  }

  public void addNamedDataSchemas(Set<String> names) {
    names.forEach(this::addNamedDataSchema);
  }

  private void addNamedDataSchema(String name) {
    NamedDataSchema schema = _parser.getSchemaResolver().existingDataSchema(name);
    if (schema == null) {
      throw new IllegalArgumentException(name + " does not exist!");
    }

    Identifier identifier = new Identifier(schema);
    if (isKnownType(identifier)) {
      return;
    }

    Path sourceFile = resolveSourceFile(schema);
    switch (schema.getType()) {
      case RECORD:
        parseDataType((RecordDataSchema) schema, sourceFile);
        break;
      case TYPEREF:
        fromDataSchema(schema, null, null, null);
        break;
      case ENUM:
        parseDataType((EnumDataSchema) schema, sourceFile);
        break;
      case FIXED:
        parseDataType((FixedDataSchema) schema, sourceFile);
        break;
      default:
        throw new IllegalArgumentException("Don't know what to do with: " + schema);
    }

    extractDataTypes(new ResourceSchema().setSchema(name));
  }

  public Set<DataType> getDataTypes() {
    return new HashSet<>(_dataTypes.values());
  }

  public RestliType parseFromRestSpec(String schema) {
    return fromDataSchema(RestSpecCodec.textToSchema(schema, _parser.getSchemaResolver()), null, null, null);
  }

  public PathKey collectionPathKey(String resourceName, String resourceNamespace, CollectionSchema collection,
      Path specFile) {
    RestliType pkType;
    if (collection.getIdentifier().hasParams()) {
      RestliType keyType = parseFromRestSpec(collection.getIdentifier().getType());
      Preconditions.checkNotNull(keyType._reference, "Complex key \"key\" must be a record type");
      RestliType paramsType = parseFromRestSpec(collection.getIdentifier().getParams());
      Preconditions.checkNotNull(paramsType._reference, "Complex key \"params\" must be a record type");
      ComplexKey key =
          new ComplexKey(resourceName, resourceNamespace, specFile, keyType._reference, paramsType._reference);
      registerDataType(new DataType(key));
      pkType = key.restliType();
    } else {
      pkType = parseFromRestSpec(collection.getIdentifier().getType());
    }

    return new PathKey(collection.getIdentifier().getName(), pkType);
  }

  private void parseDataType(RecordDataSchema schema, Path sourceFile) {
    if (_rawRecords.contains(schema.getFullName())) {
      return;
    }

    // Add all included records to the parsed types
    for (NamedDataSchema included : schema.getInclude()) {
      addNamedDataSchema(included.getFullName());
    }

    List<Field> fields = new ArrayList<>();
    for (RecordDataSchema.Field field : schema.getFields()) {
      if (schema.isFieldFromIncludes(field)) {
        // Ignore included fields, they will be resolved during code generation
        continue;
      }

      fields.add(new Field(
          field.getName(),
          field.getDoc(),
          fromDataSchema(
              field.getType(),
              schema.getNamespace(),
              sourceFile,
              Arrays.asList(schema.getName(), field.getName())
          ),
          field.getOptional(),
          field.getDefault()));
    }

    List<Identifier> includes = schema.getInclude().stream()
        .map(Identifier::new)
        .collect(Collectors.toList());

    registerDataType(new DataType(new Record(schema, sourceFile, includes, fields)));
  }

  private void parseDataType(EnumDataSchema schema, Path sourceFile) {
    registerDataType(new DataType(new Enum(schema, sourceFile, schema.getSymbols(), schema.getSymbolDocs())));
  }

  private void parseDataType(FixedDataSchema schema, Path sourceFile) {
    registerDataType(new DataType(new Fixed(schema, sourceFile, schema.getSize())));
  }

  private RestliType fromDataSchema(DataSchema schema, @Nullable String namespace, @Nullable Path sourceFile,
      @Nullable List<String> hierarchy) {
    if (schema.isPrimitive()) {
      return new RestliType(primitiveType(schema));
    }

    if (schema instanceof NamedDataSchema && _rawRecords.contains(((NamedDataSchema) schema).getFullName())) {
      if (schema.getType() != RECORD) {
        throw new IllegalArgumentException(String.format(
            "Only records can be parsed as raw records (got %s for %s)",
            schema.getType(), ((NamedDataSchema) schema).getFullName()));
      }
      return RAW_RECORD;
    }

    switch (schema.getType()) {
      case TYPEREF:
        TyperefDataSchema typerefSchema = (TyperefDataSchema) schema;

        DataSchema ref = typerefSchema.getRef();
        if (ref.isPrimitive()) {
          Typeref typeref = new Typeref(typerefSchema, resolveSourceFile(typerefSchema), primitiveType(ref));
          registerDataType(new DataType(typeref));
          return new RestliType(typeref.getIdentifier());
        } else {
          if (ref.getType() != UNION) {
            // Typerefs that represent unions need corresponding hard types generated. Any other type will simply be
            // resolved and referenced directly as they will have corresponding hard types.
            _anonymousTypeRefs.add(typerefSchema.getFullName());
          }
          if (ref.getType() == TYPEREF) {
            return fromDataSchema(ref, null, null, null);
          } else {
            return fromDataSchema(
                ref,
                typerefSchema.getNamespace(),
                resolveSourceFile(typerefSchema),
                Collections.singletonList(typerefSchema.getName()));
          }
        }
      case RECORD:
      case FIXED:
      case ENUM:
        NamedDataSchema namedSchema = (NamedDataSchema) schema;
        return new RestliType(new Identifier(namedSchema.getNamespace(), namedSchema.getName()));
      case ARRAY:
        RestliType arrayType = fromDataSchema(
            ((ArrayDataSchema) schema).getItems(),
            namespace,
            sourceFile,
            Utils.append(hierarchy, "Array"));
        return new RestliType(arrayType, null);
      case MAP:
        RestliType mapType = fromDataSchema(
            ((MapDataSchema) schema).getValues(),
            namespace,
            sourceFile,
            Utils.append(hierarchy, "Map"));
        return new RestliType(null, mapType);
      case UNION:
        Preconditions.checkState(namespace != null && hierarchy != null,
            "Raw unions not supported outside of records or typerefs");
        UnionDataSchema unionSchema = (UnionDataSchema) schema;

        String name = hierarchy.stream()
            .map(Utils::exportedIdentifier)
            .collect(Collectors.joining("_"));

        List<UnionMember> members = unionSchema.getMembers().stream()
            .map(m -> m.getType().getType() == NULL
                ? null
                : new UnionMember(fromDataSchema(m.getType(), namespace, sourceFile, hierarchy), m.getUnionMemberKey()))
            .collect(Collectors.toList());

        StandaloneUnion union = new StandaloneUnion(name, namespace, sourceFile, members);
        registerDataType(new DataType(union));

        return new RestliType(union.getIdentifier());
      default:
        throw new UnknownTypeException(schema.getType());
    }
  }

  private GoPrimitive primitiveType(DataSchema schema) {
    GoPrimitive primitive = JAVA_TO_GO_PRIMTIIVE_TYPE.get(schema.getType());
    Preconditions.checkArgument(schema.isPrimitive() && primitive != null, "Unknown primitive type %s", schema);
    return primitive;
  }

  private void registerDataType(DataType type) {
    if (!RESERVED_TYPES.contains(type.getIdentifier())) {
      _dataTypes.putIfAbsent(type.getIdentifier(), type);
    }
  }

  private Path resolveSourceFile(NamedDataSchema namedSchema) {
    DataSchemaLocation location =
        _parser.getSchemaResolver().nameToDataSchemaLocations().get(namedSchema.getFullName());
    Preconditions.checkNotNull(location, "Could not resolve original location for %s", namedSchema.getFullName());
    return location.getSourceFile().toPath().toAbsolutePath();
  }

  private boolean isKnownType(Identifier identifier) {
    return _dataTypes.containsKey(identifier)
        || RESERVED_TYPES.contains(identifier)
        || _rawRecords.contains(identifier.toString())
        || _anonymousTypeRefs.contains(identifier.toString());
  }

  public static class UnknownTypeException extends RuntimeException {
    private UnknownTypeException(DataSchema.Type type) {
      super("Unknown type: " + type);
    }
  }
}
