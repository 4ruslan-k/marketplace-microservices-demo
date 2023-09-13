
import React, { useContext } from 'react';
import TextField from '@mui/material/TextField';
import * as yup from 'yup';
import { useFormik } from 'formik';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import LoadingButton from '@mui/lab/LoadingButton';
import { useFetchUser, useSendMfaOtpCode } from '../requests/userHooks';
import { SnackbarContext } from '../context/SnackbarContext';
import { mfaCodeValidation } from '../constants/validations';

interface MfaModalProps {
  handleClose: () => void;
  open: boolean;
  passwordVerificationToken: string;
}

function MfaModal({ handleClose, open, passwordVerificationToken }: MfaModalProps) {
  const { openSuccessSnackbar, openErrorSnackbar }: any = useContext(SnackbarContext);

  const { refetchUser } = useFetchUser();
  const router = useRouter();
  const { sendMfaOtpCode, isLoadingSendingOtpCode } = useSendMfaOtpCode({
    onSuccess: () => {
      handleClose();
      openSuccessSnackbar('Login successful');
      router.push('/');
      refetchUser();
    },
    onError: (error:any) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });
  const mfaSchema = yup.object().shape({
    otp: mfaCodeValidation,
  });

  const { errors, handleChange, touched, handleBlur, isValid, handleSubmit } = useFormik({
    initialValues: {
      otp: '',
    },
    validationSchema: mfaSchema,
    onSubmit: (values) => sendMfaOtpCode({ ...values, passwordVerificationToken }),
  });
  const otpError = errors.otp && touched.otp ? errors.otp : null;

  return (
    <Dialog open={open} onClose={handleClose}>
      <form onSubmit={handleSubmit} className="flex flex-col justify-center items-center">
        <DialogTitle>MFA</DialogTitle>
        <DialogContent>
          <DialogContentText>To log in, please enter your MFA OTP code</DialogContentText>
          <TextField
            fullWidth
            autoFocus
            type="text"
            variant="standard"
            label="OTP code"
            name="otp"
            id="otp"
            onChange={handleChange}
            onBlur={handleBlur}
            helperText={otpError}
            error={Boolean(otpError)}
            required
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => handleClose()}>Cancel</Button>
          <LoadingButton
            variant="contained"
            color="primary"
            type="submit"
            disabled={!isValid || isLoadingSendingOtpCode}
            loading={isLoadingSendingOtpCode}
          >
            Submit
          </LoadingButton>
        </DialogActions>
      </form>
    </Dialog>
  );
}

export default MfaModal;
