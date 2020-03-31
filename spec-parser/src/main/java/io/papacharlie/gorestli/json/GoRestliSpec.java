package io.papacharlie.gorestli.json;

import com.google.common.base.Preconditions;
import java.util.HashSet;
import java.util.Set;


public class GoRestliSpec {
  public final Set<DataType> _dataTypes = new HashSet<>();
  public final Set<Resource> _resources = new HashSet<>();

  public static class DataType {
    private Enum _enum;
    private Fixed _fixed;
    private Record _record;
    private Typeref _typeref;
    private ComplexKey _complexKey;

    public DataType(Enum anEnum) {
      _enum = Preconditions.checkNotNull(anEnum);
    }

    public DataType(Fixed fixed) {
      _fixed = Preconditions.checkNotNull(fixed);
    }

    public DataType(Record record) {
      _record = Preconditions.checkNotNull(record);
    }

    public DataType(Typeref typeref) {
      _typeref = Preconditions.checkNotNull(typeref);
    }

    public DataType(ComplexKey complexKey) {
      _complexKey = Preconditions.checkNotNull(complexKey);
    }

    @Override
    public int hashCode() {
      return getNamedType().hashCode();
    }

    @Override
    public boolean equals(Object obj) {
      if (!(obj instanceof DataType)) {
        return false;
      }
      return getNamedType().equals(((DataType) obj).getNamedType());
    }

    private NamedType getNamedType() {
      if (_enum != null) {
        return _enum;
      }
      if (_fixed != null) {
        return _fixed;
      }
      if (_record != null) {
        return _record;
      }
      if (_typeref != null) {
        return _typeref;
      }
      if (_complexKey != null) {
        return _complexKey;
      }
      throw new IllegalStateException("No NamedType specified!");
    }
  }
}
