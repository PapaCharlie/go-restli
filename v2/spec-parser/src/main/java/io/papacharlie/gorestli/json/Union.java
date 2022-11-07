package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.json.RestliType.UnionMember;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;


public class Union {
  public final boolean _hasNull;
  public final List<UnionMember> _members;

  public Union(List<UnionMember> members) {
    _hasNull = members.contains(null);
    _members = members.stream()
        .filter(Objects::nonNull)
        .collect(Collectors.toList());
  }
}
