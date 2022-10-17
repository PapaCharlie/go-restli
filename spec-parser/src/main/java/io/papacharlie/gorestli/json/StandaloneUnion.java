package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.json.RestliType.UnionMember;
import java.io.File;
import java.nio.file.Path;
import java.util.List;


public class StandaloneUnion extends NamedType {
  public final Union _union;

  public StandaloneUnion(String name, String namespace, Path sourceFile, List<UnionMember> members) {
    super(name, namespace, "", sourceFile);
    _union = new Union(members);
  }
}
