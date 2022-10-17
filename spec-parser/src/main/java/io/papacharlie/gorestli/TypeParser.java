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
import com.linkedin.restli.restspec.RestSpecCodec;
import io.papacharlie.gorestli.json.ComplexKey;
import io.papacharlie.gorestli.json.Enum;
import io.papacharlie.gorestli.json.Fixed;
import io.papacharlie.gorestli.json.GoRestliManifest.DataType;
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
import static io.papacharlie.gorestli.json.RestliType.GoPrimitive;
import static io.papacharlie.gorestli.json.RestliType.Identifier;
import static io.papacharlie.gorestli.json.RestliType.JAVA_TO_GO_PRIMTIIVE_TYPE;
import static io.papacharlie.gorestli.json.RestliType.RAW_RECORD;
import static io.papacharlie.gorestli.json.RestliType.UnionMember;


public class TypeParser {
  private final DataSchemaParser _parser;
  private final Set<String> _rawRecords;
  private final Map<Identifier, DataType> _dataTypes = new HashMap<>();

  // These types are manually defined in go-restli and should never be re-generated
  private static final Set<Identifier> RESERVED_TYPES = ImmutableSet.of(
      MethodParser.PAGING_CONTEXT,
      new Identifier("com.linkedin.restli.common", "EmptyRecord"),
      new Identifier("restlidata", "RawRecord")
  );

  public TypeParser(DataSchemaParser parser, Set<String> rawRecords) {
    _parser = parser;
    _rawRecords = rawRecords;
  }

  public void addNamedDataSchemas(Map<String, NamedDataSchema> schemas) {
    Set<String> keysRemaining = new HashSet<>(schemas.keySet());
    // If there are inter-schema dependencies, it is necessary to consume the schemas in multiple passes since some
    // schemas may not immediately be resolvable.
    while (!keysRemaining.isEmpty()) {
      Set<String> before = new HashSet<>(keysRemaining);
      keysRemaining.removeIf(key -> addNamedDataSchema(schemas.get(key)));
      if (before.equals(keysRemaining)) {
        // If the set of remaining keys does not change between passes then they will never be consumed and this loop
        // would otherwise never exit.
        throw new IllegalStateException("Failed to process the following schemas: " + before);
      }
    }
  }

  private boolean addNamedDataSchema(NamedDataSchema schema) {
    Path sourceFile = resolveSourceFile(schema);
    switch (schema.getType()) {
      case RECORD:
        return parseDataType((RecordDataSchema) schema, sourceFile);
      case TYPEREF:
        fromDataSchema(schema, null, null, null);
        return true;
      case ENUM:
        parseDataType((EnumDataSchema) schema, sourceFile);
        return true;
      case FIXED:
        parseDataType((FixedDataSchema) schema, sourceFile);
        return true;
      default:
        throw new IllegalArgumentException("Don't know what to do with: " + schema);
    }
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

  private boolean parseDataType(RecordDataSchema schema, Path sourceFile) {
    if (_rawRecords.contains(schema.getFullName())) {
      return true;
    }

    List<Field> fields = new ArrayList<>();
    for (RecordDataSchema.Field field : schema.getFields()) {
      Identifier includedFrom = schema.isFieldFromIncludes(field)
          ? new Identifier(field.getRecord())
          : null;

      RestliType fieldRestliType;
      if (includedFrom != null) {
        // Special case for transitive types (such as unions) that don't exist in PDL but do exist as standalone types
        // in go-restli. Always refer to the included record's type and skip calling fromDataSchema to prevent such
        // transient types from being generated more than once. Additionally, this logic is necessary to allow the
        // parsing to be retried if a record depends on another record that has not yet been seen.
        if (!_dataTypes.containsKey(includedFrom)) {
          return false;
        }
        Record parent = (Record) _dataTypes.get(includedFrom).getNamedType();
        fieldRestliType = parent.getField(field.getName())._type;
      } else {
        fieldRestliType = fromDataSchema(
            field.getType(),
            schema.getNamespace(),
            sourceFile,
            Arrays.asList(schema.getName(), field.getName()));
      }
      boolean optional = field.getOptional();

      fields.add(new Field(
          field.getName(),
          field.getDoc(),
          fieldRestliType,
          optional,
          field.getDefault(),
          includedFrom));
    }

    registerDataType(new DataType(new Record(schema, sourceFile, fields)));
    return true;
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
            "Only records can be parsed as raw records (got %s for %s",
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
        } else if (typerefSchema.getRef().getType() == TYPEREF) {
          return fromDataSchema(typerefSchema.getRef(), null, null, null);
        } else {
          return fromDataSchema(
              typerefSchema.getRef(),
              typerefSchema.getNamespace(),
              resolveSourceFile(typerefSchema),
              Collections.singletonList(typerefSchema.getName()));
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
    return location.getSourceFile().toPath();
  }

  public static class UnknownTypeException extends RuntimeException {
    private UnknownTypeException(DataSchema.Type type) {
      super("Unknown type: " + type);
    }
  }
}
