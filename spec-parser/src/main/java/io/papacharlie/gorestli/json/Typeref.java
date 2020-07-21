package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.json.RestliType.Identifier;
import io.papacharlie.gorestli.json.RestliType.UnionMember;
import java.io.File;
import java.util.List;


public class Typeref extends NamedType {
  public final String _primitive;
  public final Identifier _reference;
  public final RestliType _array;
  public final RestliType _map;
  public final Union _union;

  private Typeref(NamedDataSchema namedDataSchema, File sourceFile, String primitive, Identifier reference,
      RestliType array, RestliType map, Union union) {
    super(namedDataSchema, sourceFile);
    _primitive = primitive;
    _reference = reference;
    _array = array;
    _map = map;
    _union = union;
  }

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, Identifier reference) {
    this(namedDataSchema, sourceFile, null, reference, null, null, null);
  }

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, String primitive) {
    this(namedDataSchema, sourceFile, primitive, null, null, null, null);
  }

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, RestliType array, RestliType map) {
    this(namedDataSchema, sourceFile, null, null, array, map, null);
  }

  public Typeref(NamedDataSchema namedDataSchema, File sourceFile, Union union) {
    this(namedDataSchema, sourceFile, null, null, null, null, union);
  }
}
