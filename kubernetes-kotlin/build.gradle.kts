plugins {
    kotlin("jvm") version "1.7.10"
    application
}

group = "myproject"
version = "0.0.1-SNAPSHOT"
description = "quickstart-k8s-kotlin"

repositories {
    maven { // The google mirror is less flaky than mavenCentral()
        url = uri("https://maven-central.storage-download.googleapis.com/maven2/")
    }
    mavenCentral()
    mavenLocal()
}

dependencies {
    implementation("com.pulumi:pulumi:(,1.0]")
    implementation("com.pulumi:kubernetes:3.19.1")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
}

application {
    mainClass.set(
        project.findProperty("mainClass") as? String ?: "${group}.MainKt"
    )
}
