package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.Utils;
import io.papacharlie.gorestli.json.RestliType.Identifier;
import java.io.File;
import java.nio.file.Path;


public class ComplexKey extends NamedType {
  private static final String DOC_FORMAT = "Complex Key for %s";

  public final Identifier _key;
  public final Identifier _params;

  public ComplexKey(String resourceName, String resourceNamespace, Path restSpecFile, Identifier key,
      Identifier params) {
    super(complexKeyTypeName(resourceName), resourceNamespace, String.format(DOC_FORMAT, resourceName), restSpecFile);
    _key = key;
    _params = params;
  }

  private static String complexKeyTypeName(String resourceName) {
    return Utils.exportedIdentifier(resourceName) + "_ComplexKey";
  }
}
