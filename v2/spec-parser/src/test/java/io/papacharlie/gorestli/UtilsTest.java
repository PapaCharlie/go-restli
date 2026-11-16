package io.papacharlie.gorestli;

import java.nio.file.Path;
import java.nio.file.Paths;
import org.testng.annotations.Test;

import static org.testng.Assert.assertEquals;


public class UtilsTest {
  private static class Holder {
    Path _sourceFile;
  }

  @Test
  public void pathIsSerializedAsString() {
    Holder holder = new Holder();
    holder._sourceFile = Paths.get("/some/dir/file.restspec.json");

    // The runtime type is the platform specific implementation (e.g. sun.nio.fs.UnixPath), not Path itself. Gson
    // must not try to reflectively serialize it, which fails on JDK 16+ due to strong encapsulation.
    assertEquals(Utils.UGLY_GSON.toJson(holder), "{\"sourceFile\":\"/some/dir/file.restspec.json\"}");
  }

  @Test
  public void pathRoundTrips() {
    Holder holder = new Holder();
    holder._sourceFile = Paths.get("/some/dir/file.restspec.json");

    Holder parsed = Utils.UGLY_GSON.fromJson(Utils.UGLY_GSON.toJson(holder), Holder.class);

    assertEquals(parsed._sourceFile, holder._sourceFile);
  }
}
