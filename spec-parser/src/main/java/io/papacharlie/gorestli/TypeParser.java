package io.papacharlie.gorestli;

import com.google.common.base.Preconditions;
import com.linkedin.data.schema.ArrayDataSchema;
import com.linkedin.data.schema.DataSchema;
import com.linkedin.data.schema.DataSchemaLocation;
import com.linkedin.data.schema.DataSchemaResolver;
import com.linkedin.data.schema.EnumDataSchema;
import com.linkedin.data.schema.FixedDataSchema;
import com.linkedin.data.schema.MapDataSchema;
import com.linkedin.data.schema.NamedDataSchema;
import com.linkedin.data.schema.RecordDataSchema;
import com.linkedin.data.schema.TyperefDataSchema;
import com.linkedin.data.schema.UnionDataSchema;
import com.linkedin.data.schema.UnionDataSchema.Member;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestSpecCodec;
import com.linkedin.restli.tools.snapshot.gen.SnapshotGenerator;
import io.papacharlie.gorestli.json.Enum;
import io.papacharlie.gorestli.json.Fixed;
import io.papacharlie.gorestli.json.GoRestliSpec.DataType;
import io.papacharlie.gorestli.json.Record;
import io.papacharlie.gorestli.json.Record.Field;
import io.papacharlie.gorestli.json.RestliType;
import io.papacharlie.gorestli.json.Typeref;
import java.io.File;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import static com.linkedin.data.schema.DataSchema.Type.*;
import static io.papacharlie.gorestli.json.RestliType.*;


public class TypeParser {
  private final DataSchemaResolver _dataSchemaResolver;

  public TypeParser(DataSchemaResolver dataSchemaResolver) {
    _dataSchemaResolver = dataSchemaResolver;
  }

  public Set<DataType> extractDataTypes(ResourceSchema resourceSchema) {
    Set<DataType> dataTypes = new HashSet<>();
    SnapshotGenerator generator = new SnapshotGenerator(resourceSchema, _dataSchemaResolver);
    for (NamedDataSchema schema : generator.generateModelList()) {
      DataSchemaLocation location = _dataSchemaResolver.nameToDataSchemaLocations().get(schema.getFullName());
      Preconditions.checkNotNull(location, "Could not resolve original location for %s", schema.getFullName());
      File sourceFile = location.getSourceFile();
      DataType dataType;
      switch (schema.getType()) {
        case RECORD:
          dataType = parseDataType((RecordDataSchema) schema, sourceFile);
          break;
        case ENUM:
          dataType = parseDataType((EnumDataSchema) schema, sourceFile);
          break;
        case TYPEREF:
          dataType = parseDataType((TyperefDataSchema) schema, sourceFile);
          break;
        case FIXED:
          dataType = parseDataType((FixedDataSchema) schema, sourceFile);
          break;
        default:
          System.err.printf("Don't know what to do with %s%n", schema);
          continue;
      }
      dataTypes.add(dataType);
    }
    return dataTypes;
  }

  private DataType parseDataType(RecordDataSchema schema, File sourceFile) {
    List<Field> fields = new ArrayList<>();
    for (RecordDataSchema.Field field : schema.getFields()) {
      DataSchema fieldType = field.getType();
      boolean optional = field.getOptional();
      if (fieldType instanceof UnionDataSchema) {
        for (Member unionMember : ((UnionDataSchema) fieldType).getMembers()) {
          if (unionMember.getType().getType() == NULL) {
            optional = true;
            break;
          }
        }
      }

      fields.add(new Field(
          field.getName(),
          field.getDoc(),
          fromDataSchema(fieldType),
          optional,
          field.getDefault()));
    }

    return new DataType(new Record(schema, sourceFile, fields));
  }

  private DataType parseDataType(EnumDataSchema schema, File sourceFile) {
    return new DataType(new Enum(schema, sourceFile, schema.getSymbols(), schema.getSymbolDocs()));
  }

  private DataType parseDataType(TyperefDataSchema schema, File sourceFile) {
    return new DataType(new Typeref(schema, sourceFile, fromDataSchema(schema.getDereferencedDataSchema())));
  }

  private DataType parseDataType(FixedDataSchema schema, File sourceFile) {
    return new DataType(new Fixed(schema, sourceFile, schema.getSize()));
  }

  public RestliType parseFromRestSpec(String schema) {
    return fromDataSchema(RestSpecCodec.textToSchema(schema, _dataSchemaResolver));
  }

  public RestliType fromDataSchema(DataSchema schema) {
    if (JAVA_TO_GO_PRIMTIIVE_TYPE.containsKey(schema.getType())) {
      return new RestliType(JAVA_TO_GO_PRIMTIIVE_TYPE.get(schema.getType()), null, null, null, null);
    }
    switch (schema.getType()) {
      case TYPEREF:
      case RECORD:
      case FIXED:
      case ENUM:
        NamedDataSchema namedSchema = (NamedDataSchema) schema;
        return new RestliType(null, new Identifier(namedSchema.getNamespace(), namedSchema.getName()), null, null,
            null);
      case ARRAY:
        return new RestliType(null, null, fromDataSchema(((ArrayDataSchema) schema).getItems()), null, null);
      case MAP:
        return new RestliType(null, null, null, fromDataSchema(((MapDataSchema) schema).getValues()), null);
      case UNION:
        UnionDataSchema union = (UnionDataSchema) schema;
        List<UnionMember> unionMembers = union.getMembers().stream()
            .filter(m -> m.getType().getType() != DataSchema.Type.NULL)
            .map(m -> new UnionMember(fromDataSchema(m.getType()), m.getUnionMemberKey()))
            .collect(Collectors.toList());
        if (unionMembers.size() == 1) {
          return unionMembers.get(0)._type;
        } else {
          return new RestliType(null, null, null, null, unionMembers);
        }
      default:
        throw new UnknownTypeException(schema.getType());
    }
  }

  public static class UnknownTypeException extends RuntimeException {
    private UnknownTypeException(DataSchema.Type type) {
      super("Unknown type: " + type);
    }
  }
}
