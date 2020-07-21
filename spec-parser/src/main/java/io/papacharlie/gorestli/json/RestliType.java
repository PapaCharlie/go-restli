package io.papacharlie.gorestli.json;

import com.google.common.collect.ImmutableMap;
import com.linkedin.data.schema.DataSchema;
import java.util.Map;


public final class RestliType {
  public static final Map<DataSchema.Type, String> JAVA_TO_GO_PRIMTIIVE_TYPE =
      ImmutableMap.<DataSchema.Type, String>builder()
          .put(DataSchema.Type.BOOLEAN, "bool")
          .put(DataSchema.Type.INT, "int32")
          .put(DataSchema.Type.LONG, "int64")
          .put(DataSchema.Type.FLOAT, "float32")
          .put(DataSchema.Type.DOUBLE, "float64")
          .put(DataSchema.Type.STRING, "string")
          .put(DataSchema.Type.BYTES, "bytes")
          .build();

  public final String _primitive;
  public final Identifier _reference;
  public final RestliType _array;
  public final RestliType _map;

  private RestliType(String primitive, Identifier reference, RestliType array, RestliType map) {
    _primitive = primitive;
    _reference = reference;
    _array = array;
    _map = map;
  }

  public RestliType(String primitive) {
    this(primitive, null, null, null);
  }

  public RestliType(Identifier reference) {
    this(null, reference, null, null);
  }

  public RestliType(RestliType array, RestliType map) {
    this(null, null, array, map);
  }

  public static class UnionMember {
    public final RestliType _type;
    public final String _alias;

    public UnionMember(RestliType type, String alias) {
      _type = type;
      _alias = alias;
    }
  }

  public static class Identifier {
    public final String _namespace;
    public final String _name;

    public Identifier(String namespace, String name) {
      _namespace = namespace;
      _name = name;
    }
  }

  public static class UnknownTypeException extends RuntimeException {
    private UnknownTypeException(DataSchema.Type type) {
      super("Unknown type: " + type);
    }
  }
}
