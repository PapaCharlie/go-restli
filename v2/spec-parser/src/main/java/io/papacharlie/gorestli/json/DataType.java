package io.papacharlie.gorestli.json;

import com.google.common.base.Preconditions;


public class DataType {
  private Enum _enum;
  private Fixed _fixed;
  private Record _record;
  private ComplexKey _complexKey;
  private StandaloneUnion _standaloneUnion;
  private Typeref _typeref;

  public DataType(Enum anEnum) {
    _enum = Preconditions.checkNotNull(anEnum);
  }

  public DataType(Fixed fixed) {
    _fixed = Preconditions.checkNotNull(fixed);
  }

  public DataType(Record record) {
    _record = Preconditions.checkNotNull(record);
  }

  public DataType(ComplexKey complexKey) {
    _complexKey = Preconditions.checkNotNull(complexKey);
  }

  public DataType(StandaloneUnion standaloneUnion) {
    _standaloneUnion = Preconditions.checkNotNull(standaloneUnion);
  }

  public DataType(Typeref typeref) {
    _typeref = Preconditions.checkNotNull(typeref);
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

  public RestliType.Identifier getIdentifier() {
    return getNamedType().getIdentifier();
  }

  public NamedType getNamedType() {
    if (_enum != null) {
      return _enum;
    }
    if (_fixed != null) {
      return _fixed;
    }
    if (_record != null) {
      return _record;
    }
    if (_complexKey != null) {
      return _complexKey;
    }
    if (_standaloneUnion != null) {
      return _standaloneUnion;
    }
    if (_typeref != null) {
      return _typeref;
    }
    throw new IllegalStateException("No NamedType specified!");
  }
}
