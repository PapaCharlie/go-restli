package io.papacharlie.gorestli.json;

import com.linkedin.restli.restspec.CollectionSchema;
import io.papacharlie.gorestli.TypeParser;
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
  public List<PathKey> _pathKeys;
  public List<Field> _params;
  public RestliType _return;

  public static class PathKey {
    public final String _name;
    public final RestliType _type;

    public PathKey(String name, RestliType type) {
      _name = name;
      _type = type;
    }

    public static PathKey forCollection(CollectionSchema collection, TypeParser typeParser) {
      return new PathKey(
          collection.getIdentifier().getName(),
          typeParser.parseFromRestSpec(collection.getIdentifier().getType()));
    }
  }
}
