package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.NativeTyperef;


public class RestliNativeTyperef extends NativeTyperef {
  public final RestliType.GoPrimitive _primitive;

  public RestliNativeTyperef(RestliType.GoPrimitive primitive, NativeTyperef ref) {
    super(ref);
    _primitive = primitive;
  }
}
