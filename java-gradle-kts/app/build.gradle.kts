plugins {
    application
}

repositories {
    maven { // The google mirror is less flaky than mavenCentral()
        url = uri("https://maven-central.storage-download.googleapis.com/maven2/")
    }
    mavenCentral()
    mavenLocal()
}

dependencies {
    implementation("com.pulumi:pulumi:(,1.0]")
}

application {
    mainClass.set(
        project.findProperty("mainClass") as? String ?: "myproject.App"
    )
}
