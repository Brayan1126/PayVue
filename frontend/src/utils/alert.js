import Swal from 'sweetalert2';

export const showSuccess = (message) => {
  Swal.fire({
    icon: 'success',
    title: '¡Éxito!',
    text: message,
    timer: 1500,
    showConfirmButton: false,
  });
};

export const showError = (message) => {
  Swal.fire({
    icon: 'error',
    title: 'Error',
    text: message,
    timer: 2000,
    showConfirmButton: false,
  });
};

export const showConfirm = async (message) => {
  const result = await Swal.fire({
    icon: 'warning',
    title: '¿Estás seguro?',
    text: message,
    showCancelButton: true,
    confirmButtonText: 'Sí, continuar',
    cancelButtonText: 'Cancelar',
  });
  return result.isConfirmed;
};