package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.io.File;
import java.util.List;
import java.util.Map;


public class Enum extends NamedType {
  public final List<String> _symbols;
  public final Map<String, String> _symbolToDoc;

  public Enum(NamedDataSchema namedDataSchema, File sourceFile, List<String> symbols,
      Map<String, String> symbolToDoc) {
    super(namedDataSchema, sourceFile);
    _symbols = symbols;
    _symbolToDoc = symbolToDoc;
  }
}
