package io.papacharlie.gorestli.json;

import com.google.common.collect.ImmutableMap;
import com.linkedin.data.schema.DataSchema;
import java.util.Map;
import java.util.Objects;

import static io.papacharlie.gorestli.json.RestliType.GoPrimitive.*;


public final class RestliType {
  public enum GoPrimitive {
    BOOL, INT32, INT64, FLOAT32, FLOAT64, STRING, BYTES;
  }

  public static final Map<DataSchema.Type, GoPrimitive> JAVA_TO_GO_PRIMTIIVE_TYPE =
      ImmutableMap.<DataSchema.Type, GoPrimitive>builder()
          .put(DataSchema.Type.BOOLEAN, BOOL)
          .put(DataSchema.Type.INT, INT32)
          .put(DataSchema.Type.LONG, INT64)
          .put(DataSchema.Type.FLOAT, FLOAT32)
          .put(DataSchema.Type.DOUBLE, FLOAT64)
          .put(DataSchema.Type.STRING, STRING)
          .put(DataSchema.Type.BYTES, BYTES)
          .build();

  public final GoPrimitive _primitive;
  public final Identifier _reference;
  public final RestliType _array;
  public final RestliType _map;
  public final Typeref _typeref;

  private RestliType(GoPrimitive primitive, Identifier reference, RestliType array, RestliType map, Typeref typeref) {
    _primitive = primitive;
    _reference = reference;
    _array = array;
    _map = map;
    _typeref = typeref;
  }

  public RestliType(GoPrimitive primitive) {
    this(primitive, null, null, null, null);
  }

  public RestliType(Identifier reference) {
    this(null, reference, null, null, null);
  }

  public RestliType(RestliType array, RestliType map) {
    this(null, null, array, map, null);
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

    @Override
    public int hashCode() {
      return Objects.hash(_namespace, _name);
    }

    @Override
    public boolean equals(Object obj) {
      if (!(obj instanceof Identifier)) {
        return false;
      }
      Identifier identifier = (Identifier) obj;
      return _namespace.equals(identifier._namespace)
          && _name.equals(identifier._name);
    }
  }

  public static class UnknownTypeException extends RuntimeException {
    private UnknownTypeException(DataSchema.Type type) {
      super("Unknown type: " + type);
    }
  }
}
