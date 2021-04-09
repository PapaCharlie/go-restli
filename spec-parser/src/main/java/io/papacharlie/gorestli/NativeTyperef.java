package io.papacharlie.gorestli;

/**
 * Data class representing a reference to a native go type.
 */
public class NativeTyperef {
  private final String _typePackage;
  private final String _typeName;
  private final String _objectFunctionsPackage;

  public NativeTyperef(String typePackage, String typeName, String objectFunctionsPackage) {
    _typePackage = typePackage;
    _typeName = typeName;
    _objectFunctionsPackage = objectFunctionsPackage;
  }

  protected NativeTyperef(NativeTyperef other) {
    this(other._typePackage, other._typeName, other._objectFunctionsPackage);
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

  @Override
  public String toString() {
    return Utils.UGLY_GSON.toJson(this);
  }
}
