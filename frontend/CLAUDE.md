# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**NexusOJ Frontend** - A modern online judge platform built with Vue 3, featuring problem solving, contests, leaderboards, and a knowledge base with integrated code editor and AI assistance.

## Essential Commands

```bash
# Start development server (port 3000)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Format code with Prettier
npm run format

# Upload to server (custom script)
npm run upload
```

## Project Architecture

### Tech Stack
- **Framework**: Vue 3.4+ with Composition API and TypeScript
- **Build Tool**: Vite 6+
- **Routing**: Vue Router with file-based page components
- **State**: Pinia with localStorage persistence
- **UI**: Naive UI component library
- **Styling**: Tailwind CSS 4+ with PostCSS
- **Icons**: Lucide Vue Next
- **Code Editor**: CodeMirror 6 with multi-language support
- **Rich Text**: Markdown-it with KaTeX for math expressions
- **Data Visualization**: ECharts with Vue wrapper
- **HTTP Client**: Axios with interceptors
- **Others** : Vueuse[useClipboard, useDebounce, useLocalStorage, useToggle]
  
### Directory Structure

\`\`
src/
├── assets/          # Theme CSS (dark.css, light.css, main.css)
├── components/      # Reusable Vue components
├── composables/     # Vue composition functions (useTheme.ts, useTitle.ts)
├── constants/       # Application constants
├── layouts/        # Layout wrappers (Layout.vue with sidebar, UserLayout.vue)
├── pages/          # Route components matching URL structure
├── router/         # Vue Router configuration with route definitions
├── services/        # API service layer (Request.ts for HTTP calls)
├── stores/         # Pinia stores for state management
├── types.ts        # TypeScript type definitions
├── utils/          # Utility functions
├── App.vue         # Root component with ThemeProvider
└── main.ts         # Application entry point
\`\`

### Key Architectural Patterns

1. **Page-Based Routing**: Each major route has a dedicated page component in \`src/pages/\`
2. **Composition API**: All components use \`<script setup lang="ts">\` pattern
3. **Type Safety**: Full TypeScript with strict mode and custom types
4. **Theme System**: CSS variables with dark/light mode switching via \`useTheme()\` composable
5. **State Persistence**: User data persisted via localStorage and Pinia

### Theme Architecture

- **CSS Variables**: All colors defined in \`assets/dark.css\` and \`assets/light.css\`
- **Scope**: Variables follow naming pattern \`--{category}-{variant}\`
- **Usage**: Components access via inline \` :style="{ color: 'var(--text-primary)' }"\`
- **Hero Section**: Uses dedicated \`--hero-*\` variables for independent styling

### Component Conventions

- **Naming**: PascalCase for components (e.g., \`ActivityChart.vue\`)
- **Props**: Explicitly typed with TypeScript interfaces
- **Events**: Emits with descriptive names
- **Slots**: Use named slots for clarity
- **Styles**: Prefer Tailwind utility classes, custom CSS in \`<style scoped>\` when needed

### Service Layer

All API calls go through \`src/services/api.ts\`:
- \`Request.get(url, config)\` - GET requests
- \`Request.post(url, data, config)\` - POST requests
- Base URL configured in service, proxied in Vite config

### State Management (Pinia)

Stores in \`src/stores/\`:
- \`user.ts\` - User authentication and profile data
- \`theme.ts\` - Theme preferences
- Persistence via \`localStorage\` and Pinia plugins

## Common Workflows

### Adding a New Page

1. Create component in \`src/pages/[PageName].vue\`
2. Add route in \`src/router/index.ts\`:
   \`\`\`typescript
   {
     path: '/pagename',
     name: 'PageName',
     component: () => import('@/pages/PageName.vue')
   }
   \`\`\`
3. Add navigation link in \`src/components/Navbar.vue\` if needed
4. Update \`src/types.ts\` if new data structures are needed

### Creating a Reusable Component

1. Create in \`src/components/[ComponentName].vue\`
2. Use TypeScript interfaces for props:
   \`\`\`typescript
   interface Props {
     data: Blog[]
     loading?: boolean
   }
   \`\`\`
3. Import and use in parent components
4. Export if used across multiple files

### Adding Theme Variables

For new theme-independent sections (like hero section):
1. Add variables to both \`assets/dark.css\` and \`assets/light.css\`
2. Use naming convention: \`--{section}-{property}\`
3. Access via inline styles: \` :style="{ color: 'var(--hero-title-color)' }"\`

## Code Editor Integration

The platform uses CodeMirror 6 with:
- **Languages**: JavaScript, Python, Java, C++
- **Execution**: WebSocket communication with backend
- **Themes**: One Dark theme for consistency
- **Configuration**: Located in parent component pages (e.g., Problem.vue)

## Important Constraints

1. **Use Tailwind Classes**: Prefer utility classes over custom CSS
   - Bad: \`class="my-custom-class"\`
   - Good: \`class="flex items-center gap-4"\`

2. **Gradient Class Names**: In Tailwind v4, use \`bg-linear-to-*\` not \`bg-gradient-to-*\`
   - Bad: \`class="bg-gradient-to-r"\`
   - Good: \`class="bg-linear-to-r"\`

3. **Arbitrary Values**: Use standard spacing scale where possible
   - Bad: \`class="min-w-[280px]"\`
   - Good: \`class="min-w-70"\`

4. **Theme Variables**: Always use CSS variables for themable values
   - Bad: \`class="text-white dark:text-gray-800"\`
   - Good: \` :style="{ color: 'var(--text-primary)' }"\`

5. **Responsive Design**: Mobile-first approach
   - Default: Mobile styles
   - Modifiers: \`md:\`, \`lg:\`, \`xl:\` for larger screens

## Testing

Currently no testing framework is configured. When adding tests:
1. Install Vitest: \`npm install -D vitest @vue/test-utils\`
2. Add test scripts to package.json
3. Create \`*.spec.ts\` or \`*.test.ts\` files alongside components
4. Test composition functions in \`src/composables/__tests__/\`

## API Integration

- Base URL: Configured in Vite proxy for development
- Endpoints: Defined in \`src/services/api.ts\`
- Request/Response: Typed via TypeScript interfaces
- Error Handling: Try-catch blocks with console logging

## Linting & Formatting

- **Prettier**: Auto-formats on save with 2-space indentation
- **TypeScript**: Strict mode enabled in \`tsconfig.json\`
- Run \`npm run format\` before committing

## Development Tips

1. **Hot Module Replacement (HMR)** works for most changes
2. **Router guards** are configured in \`src/router/index.ts\`
3. **Meta tags** are managed in \`src/composables/useTitle.ts\`
4. **Naive UI components** follow their documentation for best practices
5. **Pinia devtools** available in browser for state inspection

## Performance Considerations

1. **Lazy Loading**: Routes use dynamic imports for code splitting
2. **Image Optimization**: Use WebP or responsive images when possible
3. **Bundle Size**: Monitor \`dist/\` size after builds
4. **API Debouncing**: Implement for search/filter operations

## File Naming

- **Components**: PascalCase (e.g., \`ActivityChart.vue\`)
- **Composables**: camelCase with \`use\` prefix (e.g., \`useTheme.ts\`)
- **Types**: camelCase (e.g., \`blog.ts\` inside types/)
- **Utilities**: camelCase (e.g., \`formatDate.ts\`)
