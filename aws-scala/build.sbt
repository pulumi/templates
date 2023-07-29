val scala3Version = "3.2.2"

lazy val root = project
  .in(file("."))
  .settings(
    name := "Scala 3 Project Template",
    version := "0.1.0-SNAPSHOT",

    scalaVersion := scala3Version,

    libraryDependencies += "org.scalameta" %% "munit" % "0.7.29" % Test,
    libraryDependencies += "com.pulumi" % "pulumi" % "0.7.1",
    libraryDependencies += "com.pulumi" % "aws" % "5.28.0"
  )
