package io.papacharlie.gorestli;

import com.linkedin.restli.restspec.ActionSchema;
import com.linkedin.restli.restspec.ActionSchemaArray;
import com.linkedin.restli.restspec.CollectionSchema;
import com.linkedin.restli.restspec.FinderSchema;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestMethodSchema;
import com.linkedin.restli.restspec.SimpleSchema;
import io.papacharlie.gorestli.json.PathKey;
import io.papacharlie.gorestli.json.Resource;
import io.papacharlie.gorestli.json.ResourcePathSegment;
import io.papacharlie.gorestli.json.RestliType;
import java.nio.file.Path;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;


public class ResourceParser {
  private final ResourceSchema _schema;
  private final TypeParser _typeParser;
  private final List<String> _namespaceChain;
  private final List<ResourcePathSegment> _resourcePathSegments;
  private final Path _resourceFile;
  private final MethodParser _methodParser;

  private ResourceParser(ResourceSchema schema, Path resourceFile, TypeParser typeParser, List<String> namespaceChain,
      List<ResourcePathSegment> resourcePathSegments) {
    _schema = schema;
    _resourceFile = resourceFile;
    _typeParser = typeParser;
    _namespaceChain = Utils.append(namespaceChain, schema.getName());

    PathKey key;
    if (_schema.getCollection() != null) {
      key = _typeParser.collectionPathKey(_schema.getName(), namespace(), _schema.getCollection(), _resourceFile);
    } else {
      key = null;
    }
    _resourcePathSegments = Utils.append(resourcePathSegments, new ResourcePathSegment(schema.getName(), key));

    _methodParser = new MethodParser(_typeParser, _schema);
  }

  private ResourceParser(ResourceParser parent, ResourceSchema subResource) {
    this(
        subResource,
        parent._resourceFile,
        parent._typeParser,
        parent._namespaceChain,
        parent._resourcePathSegments);
  }

  private ResourceParser(ResourceSchema schema, Path resourceFile, TypeParser typeParser) {
    this(
        schema,
        resourceFile,
        typeParser,
        Collections.singletonList(schema.getNamespace()),
        Collections.emptyList());
  }

  private Set<Resource> parse() {
    // Skip parsing association resources and all subresources
    if (_schema.hasAssociation()) {
      return Collections.emptySet();
    }

    Resource resource = newResource();

    Set<Resource> resourcesAndSubResources = new HashSet<>();
    resourcesAndSubResources.add(resource);

    if (_schema.getActionsSet() != null) {
      addActions(resource, _schema.getActionsSet().getActions(), false);
    }

    if (_schema.getSimple() != null) {
      SimpleSchema simple = _schema.getSimple();
      addActions(resource, simple.getActions(), false);
      addRestMethods(resource, simple.getMethods());

      for (ResourceSchema subResource : Utils.emptyIfNull(simple.getEntity().getSubresources())) {
        resourcesAndSubResources.addAll(new ResourceParser(this, subResource).parse());
      }
    }

    if (_schema.getCollection() != null) {
      CollectionSchema collection = _schema.getCollection();
      addActions(resource, collection.getActions(), false);
      addActions(resource, collection.getEntity().getActions(), true);
      addRestMethods(resource, collection.getMethods());

      for (FinderSchema finder : Utils.emptyIfNull(collection.getFinders())) {
        resource.addMethod(_methodParser.newFinderMethod(finder));
      }

      for (ResourceSchema subResource : Utils.emptyIfNull(collection.getEntity().getSubresources())) {
        resourcesAndSubResources.addAll(new ResourceParser(this, subResource).parse());
      }
    }

    return resourcesAndSubResources;
  }

  private Resource newResource() {
    RestliType resourceType = _schema.hasSchema()
        ? _typeParser.parseFromRestSpec(_schema.getSchema())
        : null;
    Set<String> createOnlyFields = Utils.extractAnnotationValues(_schema, "createOnly");
    Set<String> readOnlyFields = Utils.extractAnnotationValues(_schema, "readOnly");
    return new Resource(
        namespace(),
        _schema.getDoc(),
        _resourceFile,
        resourceType,
        readOnlyFields,
        createOnlyFields,
        _resourcePathSegments
    );
  }

  private String namespace() {
    return String.join(".", _namespaceChain);
  }

  private void addRestMethods(Resource resource, List<RestMethodSchema> restMethods) {
    for (RestMethodSchema restMethod : Utils.emptyIfNull(restMethods)) {
      resource.addMethod(_methodParser.newRestMethod(restMethod));
    }
  }

  private void addActions(Resource resource, ActionSchemaArray actions, boolean onEntity) {
    for (ActionSchema action : Utils.emptyIfNull(actions)) {
      resource.addMethod(_methodParser.newActionMethod(action, onEntity));
    }
  }

  public static Set<Resource> parse(ResourceSchema schema, Path resourceFile, TypeParser typeParser) {
    return new ResourceParser(schema, resourceFile, typeParser).parse();
  }
}
