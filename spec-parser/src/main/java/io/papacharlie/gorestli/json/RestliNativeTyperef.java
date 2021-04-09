package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.NativeTyperef;


public class RestliNativeTyperef extends NativeTyperef {
  public final String _originalTypeName;
  public final RestliType.GoPrimitive _primitive;

  public RestliNativeTyperef(String originalTypeName, RestliType.GoPrimitive primitive, NativeTyperef ref) {
    super(ref);
    _originalTypeName = originalTypeName;
    _primitive = primitive;
  }
}
