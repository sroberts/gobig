# Mermaid Diagram Test Suite

Testing various Mermaid diagram types in gobig

---

## Flowchart

```mermaid
graph TD
    A[Start] --> B{Decision?}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
    C --> E[End]
    D --> E
```

---

## Sequence Diagram

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Server
    User->>Browser: Click Link
    Browser->>Server: HTTP Request
    Server-->>Browser: HTML Response
    Browser-->>User: Display Page
```

---

## Class Diagram

```mermaid
classDiagram
    class Animal {
        +String name
        +int age
        +makeSound()
    }
    class Dog {
        +String breed
        +bark()
    }
    class Cat {
        +String color
        +meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
```

---

## State Diagram

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Processing: Start
    Processing --> Success: Complete
    Processing --> Error: Fail
    Success --> [*]
    Error --> Idle: Retry
```

---

## Entity Relationship Diagram

```mermaid
erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
    CUSTOMER {
        string name
        string email
    }
    ORDER {
        int orderNumber
        date orderDate
    }
    LINE-ITEM {
        string productCode
        int quantity
    }
```

---

## Gantt Chart

```mermaid
gantt
    title Project Timeline
    dateFormat  YYYY-MM-DD
    section Planning
    Research           :a1, 2024-01-01, 30d
    Design             :a2, after a1, 20d
    section Development
    Backend            :a3, 2024-02-01, 45d
    Frontend           :a4, after a3, 30d
    section Testing
    QA Testing         :a5, after a4, 15d
```

---

## Pie Chart

```mermaid
pie title Programming Language Usage
    "JavaScript" : 35
    "Python" : 25
    "Go" : 20
    "Rust" : 12
    "Other" : 8
```

---

## Git Graph

```mermaid
gitGraph
    commit id: "Initial commit"
    branch develop
    checkout develop
    commit id: "Add feature"
    checkout main
    merge develop
    commit id: "Release v1.0"
```

---

## Journey Diagram

```mermaid
journey
    title User Shopping Experience
    section Browse
      Visit website: 5: User
      Search products: 4: User
    section Purchase
      Add to cart: 3: User
      Checkout: 2: User
      Payment: 1: User
    section Post-Purchase
      Receive confirmation: 5: User
      Track shipment: 4: User
```

---

## Mindmap

```mermaid
mindmap
  root((gobig))
    Features
      Markdown Support
      Layouts
      Themes
      Mermaid Diagrams
    Advantages
      Simple
      Fast
      Single File
    Use Cases
      Presentations
      Documentation
      Teaching
```

---

## Timeline

```mermaid
timeline
    title History of Web Technologies
    1991 : HTML created
    1995 : JavaScript released
    2006 : jQuery launched
    2010 : AngularJS introduced
    2013 : React released
    2020 : Modern frameworks mature
```

---

## Requirement Diagram

```mermaid
requirementDiagram
    requirement web_req {
        id: 1
        text: System shall support web interface
        risk: medium
        verifymethod: test
    }

    requirement api_req {
        id: 2
        text: System shall provide REST API
        risk: low
        verifymethod: inspection
    }

    web_req - contains -> api_req
```

---

# Summary

All Mermaid diagram types tested with **dark theme**

Each diagram should display with:
- **Light text** on dark backgrounds
- **Visible borders** and arrows
- **Proper contrast** for readability
