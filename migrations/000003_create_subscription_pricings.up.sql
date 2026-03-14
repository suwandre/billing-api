CREATE TABLE subscription_pricings {
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  type SMALLINT NOT NULL,
  price DECIMAL(10, 2) NOT NULL
}