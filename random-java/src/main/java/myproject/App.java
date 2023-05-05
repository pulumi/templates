package myproject;

import com.pulumi.Pulumi;
import com.pulumi.random.RandomPet;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var username = new RandomPet("username");

            ctx.export("name", username.id());
        });
    }
}