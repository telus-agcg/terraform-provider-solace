package com.telusagriculture.terraform_provider_generator;

import org.openapitools.codegen.*;
import org.openapitools.codegen.languages.AbstractGoCodegen;
import org.openapitools.codegen.meta.features.SchemaSupportFeature;

import com.google.common.collect.Iterables;
import com.samskivert.mustache.Mustache;
import com.samskivert.mustache.Template.Fragment;
import org.openapitools.codegen.model.ModelsMap;

import java.util.*;
import java.io.IOException;
import java.io.Writer;

public class TerraformProviderGenerator extends AbstractGoCodegen {

  // source folder where to write the files
  protected String sourceFolder = "src";
  protected String apiVersion = "1.0.0";

  /**
   * Configures the type of generator.
   *
   * @return  the CodegenType for this generator
   * @see     org.openapitools.codegen.CodegenType
   */
  public CodegenType getTag() {
    return CodegenType.OTHER;
  }

  /**
   * Configures a friendly name for the generator.  This will be used by the generator
   * to select the library with the -g flag.
   *
   * @return the friendly name for the generator
   */
  public String getName() {
    return "terraform-provider";
  }

  /**
   * Provides an opportunity to inspect and modify operation data before the code is generated.
   */
  @SuppressWarnings("unchecked")
  @Override
  public ModelsMap postProcessModels(ModelsMap objs) {
        // The superclass determines the list of required golang imports. The actual list of imports
        // depends on which types are used, some of which are changed in the code below (but then preserved
        // and used through x-go-base-type in templates). So super.postProcessModels
        // must be invoked at the beginning of this method.
        objs = super.postProcessModels(objs);

        List<Map<String, String>> imports = (List<Map<String, String>>) objs.get("imports");

        List<Map<String, Object>> models = (List<Map<String, Object>>) objs.get("models");
        if(models == null) return objs;

        for (Map<String, Object> m : models) {
            Object v = m.get("model");
            if (v instanceof CodegenModel) {
                CodegenModel model = (CodegenModel) v;
                if (model.isEnum) {
                    continue;
                }

                for (CodegenProperty param : Iterables.concat(model.vars, model.allVars, model.requiredVars, model.optionalVars)) {
                    param.vendorExtensions.put("x-go-base-type", param.dataType);

                    if(param.minLength != null || param.maxLength != null || param.allowableValues != null) {
                        imports.add(createMapping("import", "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"));
                    }

                    if(param.pattern != null) {
                        imports.add(createMapping("import", "regexp"));
                    }

                    if (!param.isNullable || param.isContainer || param.isFreeFormObject
                            || (param.isAnyType && !param.isModel)) {
                        continue;
                    }
                    if (param.isDateTime) {
                        // Note this could have been done by adding the following line in processOpts(),
                        // however, we only want to represent the DateTime object as NullableTime if
                        // it's marked as nullable in the spec.
                        //    typeMapping.put("DateTime", "NullableTime");
                        param.dataType = "NullableTime";
                    } else {
                        param.dataType = "Nullable" + Character.toUpperCase(param.dataType.charAt(0))
                                + param.dataType.substring(1);
                    }
                }
            }
        }

        return objs;
  }

  /**
   * Returns human-friendly help for the generator.  Provide the consumer with help
   * tips, parameters here
   *
   * @return A string value for the help message
   */
  public String getHelp() {
    return "Generates a terraform provider for this API";
  }

  public TerraformProviderGenerator() {
    super();

    modifyFeatureSet(features -> features
      .excludeSchemaSupportFeatures(SchemaSupportFeature.Polymorphism)
    );

    tfTypeMap.clear();
    tfTypeMap.put("string", "types.StringType");
    tfTypeMap.put("bool", "types.BoolType");
    tfTypeMap.put("int32", "types.Int64Type");
    tfTypeMap.put("int64", "types.Int64Type");

    // set the output folder here
    outputFolder = "generated-code/terraform-provider";

    /**
     * Models.  You can write model files using the modelTemplateFiles map.
     * if you want to create one template for file, you can do so here.
     * for multiple files for model, just put another entry in the `modelTemplateFiles` with
     * a different extension
     */
    modelTemplateFiles.put(
      "model.mustache", // the template to use
      ".go");       // the extension for each file to write

    /**
     * Api classes.  You can write classes for each Api file with the apiTemplateFiles map.
     * as with models, add multiple entries with different extensions for multiple files per
     * class
     */
    apiTemplateFiles.put(
      "api.mustache",   // the template to use
      ".sample");       // the extension for each file to write

    /**
     * Template Location.  This is the location which templates will be read from.  The generator
     * will use the resource stream to attempt to read the templates.
     */
    templateDir = "terraform-provider";

    /**
     * Api Package.  Optional, if needed, this can be used in templates
     */
    apiPackage = "api";

    /**
     * Model Package.  Optional, if needed, this can be used in templates
     */
    modelPackage = "model";

    /**
     * Additional Properties.  These values can be passed to the templates and
     * are available in models, apis, and supporting files
     */
    additionalProperties.put("apiVersion", apiVersion);
    additionalProperties.put("tfType", toTerraformTypeLambda);

    /**
     * Supporting Files.  You can write single files for the generator with the
     * entire object tree available.  If the input file has a suffix of `.mustache
     * it will be processed by the template engine.  Otherwise, it will be copied
     */
    supportingFiles.add(new SupportingFile("myFile.mustache",   // the input template or file
      "",                                                       // the destination folder, relative `outputFolder`
      "myFile.sample")                                          // the output file
    );

    supportingFiles.add(new SupportingFile("validation.go.mustache", "", "validation.go"));

    /**
     * Language Specific Primitives.  These types will not trigger imports by
     * the client generator
     */
    languageSpecificPrimitives = new HashSet<String>(
      Arrays.asList(
        "Type1",      // replace these with your types
        "Type2")
    );
  }

  private Map<String, String> tfTypeMap = new HashMap<>();
  public String toTerraformType(String type) {
    return tfTypeMap.getOrDefault(type, type + "Type");
  }
  private Mustache.Lambda toTerraformTypeLambda = new Mustache.Lambda() {
    @Override
    public void execute(Fragment frag, Writer out) throws IOException {
      out.write(toTerraformType(frag.execute()));
    }
  };

  /**
   * Escapes a reserved word as defined in the `reservedWords` array. Handle escaping
   * those terms here.  This logic is only called if a variable matches the reserved words
   *
   * @return the escaped term
   */
  @Override
  public String escapeReservedWord(String name) {
    return "_" + name;  // add an underscore to the name
  }

  /**
   * Location to write model files.  You can use the modelPackage() as defined when the class is
   * instantiated
   */
  public String modelFileFolder() {
    return outputFolder;
  }

  /**
   * Location to write api files.  You can use the apiPackage() as defined when the class is
   * instantiated
   */
  @Override
  public String apiFileFolder() {
    return outputFolder;
  }

  /**
   * override with any special text escaping logic to handle unsafe
   * characters so as to avoid code injection
   *
   * @param input String to be cleaned up
   * @return string with unsafe characters removed or escaped
   */
  @Override
  public String escapeUnsafeCharacters(String input) {
    //TODO: check that this logic is safe to escape unsafe characters to avoid code injection
    return input;
  }

  /**
   * Escape single and/or double quote to avoid code injection
   *
   * @param input String to be cleaned up
   * @return string with quotation mark removed or escaped
   */
  public String escapeQuotationMark(String input) {
    //TODO: check that this logic is safe to escape quotation mark to avoid code injection
    return input.replace("\"", "\\\"");
  }
}