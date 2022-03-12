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
  public boolean _onEntity;
  public List<Field> _params;
  public RestliType _return;
  public boolean _returnEntity;
  public boolean _pagingSupported;
  public RestliType _metadata;
}
