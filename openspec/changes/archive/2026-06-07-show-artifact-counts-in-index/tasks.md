## 1. Index View — Section Title Counts

- [x] 1.1 Add item count to "Active Changes" section title in `renderIndexContent()`: show `"Active Changes (N)"` where N is `len(m.project.Changes)`, only when N > 0
- [x] 1.2 Add item count to "Specifications" section title in `renderIndexContent()`: show `"Specifications (N)"` where N is `len(m.projectSpecs)`, only when N > 0
- [x] 1.3 Add item count to "Archived Changes" section title in `renderIndexContent()`: show `"Archived Changes (N)"` where N is `len(m.index.ArchiveChanges)`, only when N > 0
- [x] 1.4 Run tests and verify the index view renders correctly with and without items in each section
