// Shared UI component types

export interface BaseComponentProps {
  disabled?: boolean
  loading?: boolean
}

export interface SizeProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

export interface VariantProps {
  variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'warning'
}

export type ButtonProps = BaseComponentProps & SizeProps & VariantProps
