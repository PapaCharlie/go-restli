package io.papacharlie.gorestli;

import java.util.Objects;


/**
 * Data class representing a reference to a native go type.
 */
public class NativeTyperef {
  private final String _typePackage;
  private final String _typeName;
  private final String _objectFunctionsPackage;
  private final boolean _isValidMapKey;

  public NativeTyperef(String typePackage, String typeName, String objectFunctionsPackage, boolean isValidMapKey) {
    _typePackage = typePackage;
    _typeName = typeName;
    _objectFunctionsPackage = objectFunctionsPackage;
    _isValidMapKey = isValidMapKey;
  }

  public NativeTyperef(String typePackage, String typeName, String objectFunctionsPackage) {
    this(typePackage, typeName, objectFunctionsPackage, false);
  }

  protected NativeTyperef(NativeTyperef other) {
    this(other._typePackage, other._typeName, other._objectFunctionsPackage, other._isValidMapKey);
  }

  /**
   * Returns the package that contains the native type
   */
  public String getTypePackage() {
    return _typePackage;
  }

  /**
   * Returns the name of the native type
   */
  public String getTypeName() {
    return _typeName;
  }

  /**
   * Returns the package that contains the object functions (i.e. "Unmarshal{typeName}", "Marshal{typeName}",
   * "Equals{typeName}", "ComputeHash{typeName}" and "ZeroValue{typeName}").
   */
  public String getObjectFunctionsPackage() {
    return _objectFunctionsPackage;
  }

  /**
   * Returns whether or not this type is a valid map key (aka that == and != are well defined on it)
   */
  public boolean isValidMapKey() {
    return _isValidMapKey;
  }

  @Override
  public String toString() {
    return Utils.UGLY_GSON.toJson(this);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    NativeTyperef that = (NativeTyperef) o;
    return Objects.equals(_typePackage, that._typePackage)
        && Objects.equals(_typeName, that._typeName)
        && Objects.equals(_objectFunctionsPackage, that._objectFunctionsPackage);
  }

  @Override
  public int hashCode() {
    return Objects.hash(_typePackage, _typeName, _objectFunctionsPackage);
  }
}
