I'm sharing my Go calculator package and two existing table-driven unit tests for the Add and Divide functions.

Act as a senior Go software engineer and perform a thorough review of my implementation and test coverage.

Please:

1. Identify missing edge cases for every calculator operation, including:
   - Integer and floating-point overflow.
   - Division by zero.
   - NaN and positive/negative infinity.
   - Negative zero and floating-point precision.
   - Negative operands and decimal values.
   - Invalid inputs and boundary conditions.
   - Square root of negative numbers and other operation-specific edge cases.

2. Explain which cases are relevant to my implementation and which are unnecessary or redundant. Do not introduce tests for behavior my API does not support.

3. Write table-driven unit tests for all remaining functions, following exactly the naming conventions, structure, assertion style, and conventions of my existing Add and Divide tests.

4. Use Go's standard testing package. Avoid unnecessary third-party dependencies.

5. Do not modify my production code unless you identify a concrete bug. If a change is necessary, explain why and suggest the smallest possible fix.

6. Ensure the tests are deterministic, independent, readable, and maintainable.

7. Provide the complete test code and briefly explain the purpose of each new test case.

Prioritize correctness, idiomatic Go practices, and meaningful test coverage over maximizing the number of test cases.