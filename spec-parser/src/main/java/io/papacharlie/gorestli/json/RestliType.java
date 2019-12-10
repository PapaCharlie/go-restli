package io.papacharlie.gorestli.json;

import com.google.common.collect.ImmutableMap;
import com.linkedin.data.schema.DataSchema;
import java.util.Collections;
import java.util.List;
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
  public final List<UnionMember> _union;

  public RestliType(String primitive, Identifier reference, RestliType array, RestliType map, List<UnionMember> union) {
    _primitive = primitive;
    _reference = reference;
    _array = array;
    _map = map;
    _union = (union == null) ? null : Collections.unmodifiableList(union);
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
