

```terraform

resource "spotinstadmin_account" "this" {
  name            = "${var.account_name}"
  aws_role_arn    = "${aws_iam_role.spotinst.arn}"
  aws_external_id = "${local.external_id}"
}

resource "spotinstadmin_programmatic_user" "this" {
  name        = "${var.account_name}"
  account_id  = "${spotinstadmin_account.this.id}"
  description = "Programmatic user for ${var.account_name} account"
}
```
