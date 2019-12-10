package io.papacharlie.gorestli.json;

import com.google.common.base.Preconditions;
import java.util.HashSet;
import java.util.Set;


public class GoRestliSpec {
  public final Set<DataType> _dataTypes = new HashSet<>();
  public final Set<Resource> _resources = new HashSet<>();

  public static class DataType {
    public final Enum _enum;
    public final Fixed _fixed;
    public final Record _record;
    public final Typeref _typeref;

    private DataType(Enum anEnum, Fixed fixed, Record record, Typeref typeref) {
      _enum = anEnum;
      _fixed = fixed;
      _record = record;
      _typeref = typeref;
      Preconditions.checkNotNull(getNamedType(), "Must specify at least one type");
    }

    public DataType(Enum anEnum) {
      this(anEnum, null, null, null);
    }

    public DataType(Fixed fixed) {
      this(null, fixed, null, null);
    }

    public DataType(Record record) {
      this(null, null, record, null);
    }

    public DataType(Typeref typeref) {
      this(null, null, null, typeref);
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
      return _typeref;
    }
  }
}
