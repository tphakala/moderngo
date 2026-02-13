# Changelog

## 2025-02-13

- **SliceRepeat**: Added false-positive caveat to report message. When the appended expression depends on the loop variable (flatMap pattern), the expanded message now makes this obvious. Documented as known false positive in CLAUDE.md.
- **Real-world validation**: Ran rules against vainu2 project. 11 findings: 10 correct, 1 SliceRepeat false positive (flatMap pattern). ErrorsAsType (7), NewWithExpression (1) all correct.
