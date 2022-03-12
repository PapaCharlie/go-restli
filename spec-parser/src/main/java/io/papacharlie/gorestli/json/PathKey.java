package io.papacharlie.gorestli.json;

public class PathKey {
  public final String _name;
  public final RestliType _type;

  public PathKey(String name, RestliType type) {
    _name = name;
    _type = type;
  }
}
