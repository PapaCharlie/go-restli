package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.io.File;


public class Fixed extends NamedType {
  public final int _size;

  public Fixed(NamedDataSchema namedDataSchema, File sourceFile, int size) {
    super(namedDataSchema, sourceFile);
    _size = size;
  }
}
