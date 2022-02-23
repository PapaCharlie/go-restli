package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.json.Record.Field;
import java.util.List;


public class Method {
  public enum MethodType {
    REST_METHOD,
    ACTION,
    FINDER
  }

  public MethodType _methodType;
  public String _name;
  public String _doc;
  public String _path;
  public boolean _onEntity;
  public PathKey _entityPathKey;
  public List<PathKey> _pathKeys;
  public List<Field> _params;
  public RestliType _return;
  public boolean _returnEntity;
  public boolean _pagingSupported;
  public RestliType _metadata;

  public static class PathKey {
    public final String _name;
    public final RestliType _type;

    public PathKey(String name, RestliType type) {
      _name = name;
      _type = type;
    }
  }
}
