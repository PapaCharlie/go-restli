package io.papacharlie.gorestli;

import com.google.common.base.Preconditions;
import com.google.common.collect.ImmutableSet;
import com.linkedin.restli.common.ResourceMethod;
import com.linkedin.restli.restspec.ActionSchema;
import com.linkedin.restli.restspec.CollectionSchema;
import com.linkedin.restli.restspec.FinderSchema;
import com.linkedin.restli.restspec.ParameterSchema;
import com.linkedin.restli.restspec.ParameterSchemaArray;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestMethodSchema;
import io.papacharlie.gorestli.json.Method;
import io.papacharlie.gorestli.json.Method.MethodType;
import io.papacharlie.gorestli.json.Method.PathKey;
import io.papacharlie.gorestli.json.Record.Field;
import io.papacharlie.gorestli.json.RestliType;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Set;

import static io.papacharlie.gorestli.json.Method.MethodType.*;


public class MethodParser {
  // The following specify the keys as an "ids" query parameter
  private static final Set<ResourceMethod> BATCH_METHODS_WITH_IDS_PARAM = ImmutableSet.<ResourceMethod>builder()
      .add(ResourceMethod.BATCH_GET)
      .add(ResourceMethod.BATCH_DELETE)
      .add(ResourceMethod.BATCH_UPDATE)
      .add(ResourceMethod.BATCH_PARTIAL_UPDATE)
      .build();
  private static final Set<ResourceMethod> NO_KEY_METHODS = ImmutableSet.<ResourceMethod>builder()
      .add(ResourceMethod.CREATE)
      .add(ResourceMethod.BATCH_CREATE)
      .add(ResourceMethod.GET_ALL)
      .addAll(BATCH_METHODS_WITH_IDS_PARAM)
      .build();

  private final TypeParser _typeParser;
  private final ResourceSchema _resource;
  private final RestliType _resourceSchema;
  private final String _path;
  private final List<PathKey> _pathKeys;
  private final String _entityPath;
  private final PathKey _entityPathKey;

  public MethodParser(TypeParser typeParser, ResourceSchema resource, List<PathKey> parentPathKeys,
      PathKey entityPathKey) {
    _typeParser = typeParser;
    _resource = resource;
    if (_resource.getSchema() != null) {
      _resourceSchema = _typeParser.parseFromRestSpec(_resource.getSchema());
    } else {
      _resourceSchema = null;
    }
    _path = resource.getPath();
    _pathKeys = parentPathKeys;
    if (resource.getCollection() != null && entityPathKey != null) {
      CollectionSchema collectionSchema = resource.getCollection();
      _entityPath = collectionSchema.getEntity().getPath();
      _entityPathKey = entityPathKey;
    } else {
      _entityPath = null;
      _entityPathKey = null;
    }
  }

  public Method newActionMethod(ActionSchema action, boolean isActionOnEntity) {
    Method method = newMethod(action.getName(), ACTION, isActionOnEntity);
    method._doc = action.getDoc();
    method._params = toFieldList(action.getParameters());

    if (action.getReturns() != null) {
      method._return = _typeParser.parseFromRestSpec(action.getReturns());
    }

    return method;
  }

  public Method newFinderMethod(FinderSchema finder) {
    Method method = newMethod(finder.getName(), FINDER, false);
    method._doc = finder.getDoc();
    method._params = toFieldList(finder.getParameters());
    method._return = _resourceSchema;
    return method;
  }

  public Method newRestMethod(RestMethodSchema restMethod) {
    boolean onEntity;
    if (_resource.getSimple() != null) {
      // simple resources don't have entities
      onEntity = false;
    } else {
      onEntity = !NO_KEY_METHODS.contains(ResourceMethod.fromString(restMethod.getMethod()));
    }

    Method method = newMethod(restMethod.getMethod(), REST_METHOD, onEntity);
    method._params = toFieldList(restMethod.getParameters());
    method._return = _resourceSchema;
    method._returnEntity = Utils.supportsReturnEntity(restMethod);
    return method;
  }

  private List<Field> toFieldList(ParameterSchemaArray parameters) {
    if (parameters == null || parameters.isEmpty()) {
      return Collections.emptyList();
    }
    List<Field> fields = new ArrayList<>();
    for (ParameterSchema parameter : parameters) {
      fields.add(new Field(
          parameter.getName(),
          parameter.getDoc(),
          _typeParser.parseFromRestSpec(parameter.getType()),
          (parameter.hasOptional() && parameter.isOptional()) || parameter.hasDefault()));
    }
    return fields;
  }

  private Method newMethod(String name, MethodType methodType, boolean onEntity) {
    Method method = new Method();
    method._name = name;
    method._methodType = methodType;
    method._onEntity = onEntity;
    method._entityPathKey = _entityPathKey;

    if (onEntity) {
      Preconditions.checkNotNull(_entityPathKey);
      method._path = _entityPath;
      method._pathKeys = Utils.append(_pathKeys, _entityPathKey);
    } else {
      method._path = _path;
      method._pathKeys = _pathKeys;
    }

    return method;
  }
}
