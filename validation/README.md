# Validation

## Introduction

The Validation component runs a set of rules over a set of subjects. A rule holds
one subject, decides whether that subject passes, and carries the message that a
failure reports. A validator groups rules by the name of the subject that each
set applies to, runs them, and collects one message per subject that failed.

There is nothing implicit here. A caller builds each rule, passes each subject
directly, and reads a plain string back.

## Rules

The other ports declare an abstract `Rule` that each rule extends, and each
subclass overrides `isValid`. Go has no abstract type and no method override, so
one struct holds the subject and the message, and a function decides whether the
subject passes:

```go
type CheckFunc func(subject any) bool

func NewRule(subject any, errorMessage string, check CheckFunc) *Rule
```

Each constructor below names that function, so a caller never writes one:

```go
// Right — the constructor carries the check and the message.
built := rule.NewMin(name, 2)

err := built.Validate()
if err != nil {
	return err
}
```

`IsValid` reports whether the subject passes. `Validate` reports a failure with a
`ValidationRuleFailureError` that carries the message of the rule, because Go
reports a failure with a returned error rather than a throw.

### Identity and presence

| Constructor    | Passes when                              |
| :------------- | :--------------------------------------- |
| `NewRequired`  | The subject carries a value              |
| `NewNotEmpty`  | The subject carries a value              |
| `NewIsEmpty`   | The subject carries no value             |
| `NewEqual`     | The subject is the value                 |
| `NewNotEqual`  | The subject is not the value             |
| `NewIsString`  | The subject is a string                  |
| `NewIsNumeric` | The subject is a number, or reads as one |
| `NewIsBool`    | The subject is a boolean                 |
| `NewEmail`     | The subject is an email address          |

The other ports read PHP's `empty`, which is true for an empty string, a zero, a
false, an empty list, and a null. Go has no such rule, so the component states
the same one for each of those types.

Warning: Go's `net/mail` reads an address that carries a display name, such as
`Melech <melech@example.com>`. A subject that carries one is not an address on
its own, so `NewEmail` fails it.

### String rules

Each one fails a subject that is not a string.

| Constructor     | Extra argument   | Passes when                      |
| :-------------- | :--------------- | :------------------------------- |
| `NewMin`        | `shortest int`   | The subject is long enough       |
| `NewMax`        | `longest int`    | The subject is short enough      |
| `NewContains`   | `needle string`  | The subject carries the needle   |
| `NewStartsWith` | `needle string`  | The subject starts with it       |
| `NewEndsWith`   | `needle string`  | The subject ends with it         |
| `NewAlpha`      | —                | Every character is a letter      |
| `NewLowercase`  | —                | The subject carries no uppercase |
| `NewUppercase`  | —                | The subject carries no lowercase |
| `NewRegex`      | `pattern string` | The subject matches the pattern  |

`NewMin` and `NewMax` count runes rather than bytes, so a subject that carries a
character outside ASCII counts once.

Warning: Go's regular expressions are RE2, and a pattern carries no delimiter. A
pattern that PHP writes as `/^\d+$/` is `^\d+$` here. RE2 also rejects lookahead,
lookbehind, and a backreference. A pattern that RE2 rejects fails every subject,
because a rule has no return to carry the failure.

### Integer rules

Each one fails a subject that is not a whole number.

| Constructor      | Extra argument | Passes when             |
| :--------------- | :------------- | :---------------------- |
| `NewGreaterThan` | `lowest int`   | The subject is above it |
| `NewLessThan`    | `highest int`  | The subject is below it |

## The Validator

```go
built := validator.NewValidator(map[string][]contract.RuleContract{
	"email": {rule.NewEmail(email)},
	"name":  {rule.NewRequired(name), rule.NewMin(name, 2)},
})

if !built.ValidateRules() {
	return built.GetErrorMessages()
}
```

A subject reports one message: the first rule that it fails. A later rule that
the same subject fails does not replace it, because a person reads one message
per field. Each message reads `subject: text`.

Go's map has no order, and the first message that a caller reads must be the same
on every run. The validator therefore walks its subjects in the order of their
names, and `GetFirstErrorMessage` returns the message of the first one that
failed.

`ValidateRules` drops the messages of the run before it, so a validator that a
caller reuses reports the run that it just made.

## Service Registration

`ValidationServiceProvider` publishes one binding key:

| Binding key                                       | Holds                    |
| :------------------------------------------------ | :----------------------- |
| `valkyrja.validation.validator.ValidatorContract` | A validator with no rule |

A caller states the rules of one validation, so the validator is bound with none.
The caller sets them with `SetRules` before it validates.
