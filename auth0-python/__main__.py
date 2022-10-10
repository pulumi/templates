import pulumi
import pulumi_auth0 as auth0

client = auth0.Client("client",
                      allowed_clients=["https://allowed.example.com"],
                      allowed_logout_urls=["https://example.com"],
                      allowed_origins=["https://example.com"],
                      app_type="regular_web",
                      callbacks=["https://example.com/callback"])

pulumi.export('client_id', client.client_id)
pulumi.export('client_secret', client.client_secret)
