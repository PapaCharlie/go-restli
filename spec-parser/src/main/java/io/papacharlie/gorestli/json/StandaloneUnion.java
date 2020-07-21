package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.json.RestliType.UnionMember;
import java.io.File;
import java.util.List;


public class StandaloneUnion extends NamedType {
  public final Union _union;

  public StandaloneUnion(String name, String namespace, File sourceFile, List<UnionMember> members) {
    super(name, namespace, "", sourceFile);
    _union = new Union(members);
  }
}
