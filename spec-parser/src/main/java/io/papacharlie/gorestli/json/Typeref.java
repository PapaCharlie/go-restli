package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.io.File;


public class Typeref extends NamedType {
  public final RestliType _ref;

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, RestliType refType) {
    super(namedDataSchema, sourceFile);
    _ref = refType;
  }
}
