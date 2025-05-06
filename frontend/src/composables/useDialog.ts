import { reactive } from 'vue'

type DialogType = 'confirm' | 'error'

interface DialogOptions {
  title: string
  message: string
  type?: DialogType
  confirmText?: string
  cancelText?: string
  onConfirm?: () => void
  onCancel?: () => void
}

export function useDialog() {
  const dialog = reactive({
    visible: false,
    title: '',
    message: '',
    type: 'confirm' as DialogType,
    confirmText: 'OK',
    cancelText: 'Cancel',
    onConfirm: () => {
      dialog.visible = false
    },
    onCancel: () => {
      dialog.visible = false
    },
  })

  function showDialog(options: DialogOptions) {
    dialog.title = options.title
    dialog.message = options.message
    dialog.type = options.type || 'confirm'
    dialog.confirmText = options.confirmText || 'OK'
    dialog.cancelText = options.cancelText || 'Cancel'

    dialog.onConfirm = () => {
      dialog.visible = false
      options.onConfirm?.()
    }

    dialog.onCancel = () => {
      dialog.visible = false
      options.onCancel?.()
    }

    dialog.visible = true
  }

  return { dialog, showDialog }
}
