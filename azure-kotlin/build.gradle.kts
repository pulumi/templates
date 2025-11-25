plugins {
    kotlin("jvm") version "2.2.20"

    application
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("com.pulumi:pulumi:1.16.3")
    implementation("com.pulumi:azure-native:3.10.1")}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    mainClass = "org.example.AppKt"
}
