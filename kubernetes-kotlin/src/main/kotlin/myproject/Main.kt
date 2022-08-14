package myproject

import com.pulumi.Context
import com.pulumi.Pulumi
import com.pulumi.kubernetes.meta_v1.outputs.ObjectMeta
import java.util.*

fun run(ctx: Context) {
    val deployment = getDeployment(mapOf("app" to "nginx"))
    val name = deployment.metadata()
        .applyValue { m: Optional<ObjectMeta> -> m.orElseThrow().name().orElse("") }
    ctx.export("name", name)
}

fun main(args: Array<String>) {
    Pulumi.run(::run)
}

