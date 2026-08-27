# Instruction

## Todo List 需求文档&设计规范

构建一个简单的, 管理待办事项的应用。它基于 sqlite3 数据库，使用 golang/gin/gorm 作为后端，使用 Vue 3.0/TpyeScript/Vite/TailWind/Shadcn 作为前端，前端按照 apple website 风格，think ultra hard, 优化 UI 和 UX，使用 Makefile 管理应用的生命周期。基于 REST Client 测试应用的 REST API。无需用户系统，当前用户可以：

- 创建/删除/完成/取消完成 待办事项

按照这个想法，帮我生成详细需求和设计文档，放在 `./spec/0001-prd-spec.md` 文件中，输出为中文。


## Implementation plan 实施计划

按照 ./spec/0001-prd-spec.md 文件中的需求和设计文档，生成一个详细的、阶段性的实施计划，放在./spec/0002-implementation-plan.md 文件中，输出为中文。

## Milestone 里程碑

按照 ./spec/0002-implementation-plan.md, 在 ./todoList/ 目录下 完整实现这个项目的 M0 阶段的所有任务。

## Testing 测试
帮我根据 REST Client 撰写一个 test.rest 文件，里面包含对所有支持的 API 的测试。

