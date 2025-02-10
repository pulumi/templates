//DEPS com.pulumi:pulumi:1.+

import com.pulumi.Pulumi;
import com.pulumi.core.Output;

public class main {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            ctx.export("exampleOutput", Output.of("example"));
        });
    }
}
