package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.io.File;
import java.nio.file.Path;


public class Fixed extends NamedType {
  public final int _size;

  public Fixed(NamedDataSchema namedDataSchema, Path sourceFile, int size) {
    super(namedDataSchema, sourceFile);
    _size = size;
  }
}
