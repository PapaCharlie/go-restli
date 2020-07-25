package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.json.RestliType.GoPrimitive;
import java.io.File;


public class Typeref extends NamedType {
  public final GoPrimitive _type;

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, GoPrimitive type) {
    super(namedDataSchema, sourceFile);
    _type = type;
  }
}
