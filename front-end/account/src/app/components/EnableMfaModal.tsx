import React, { useContext } from 'react';
import TextField from '@mui/material/TextField';
import * as yup from 'yup';
import { useFormik } from 'formik';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import LoadingButton from '@mui/lab/LoadingButton';
import { useEnableTotpMfa, useFetchUser } from '../requests/userHooks';
import { SnackbarContext } from '../context/SnackbarContext';
import { mfaCodeValidation } from '../constants/validations';

interface EnableMfaModalProps {
  handleClose: () => void;
  open: boolean;
  totpImage: string;
}

function EnableMfaModal({ handleClose, open, totpImage }: EnableMfaModalProps) {
  const { openSuccessSnackbar, openErrorSnackbar }: any = useContext(SnackbarContext);

  const { refetchUser } = useFetchUser();
  const { isLoadingEnableTotpMfa, enableTotpMfa } = useEnableTotpMfa({
    onSuccess: () => {
      handleClose();
      openSuccessSnackbar('MFA has been enabled');
      refetchUser();
    },
    onError: (error) => {
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
    onSubmit: (values) => enableTotpMfa({ code: values.otp }),
  });
  const otpError = errors.otp && touched.otp ? errors.otp : null;

  return (
    <Dialog open={open} onClose={handleClose}>
      <form onSubmit={handleSubmit} className="flex flex-col justify-center items-center">
        <DialogTitle>MFA</DialogTitle>
        <DialogContent>
          <img src={`data:image/jpeg;base64,${totpImage}`} alt="qrCode" />
          <DialogContentText>To enable MFA, please enter your MFA OTP code using the provided QR code</DialogContentText>
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
            disabled={!isValid || isLoadingEnableTotpMfa}
            loading={isLoadingEnableTotpMfa}
          >
            Submit
          </LoadingButton>
        </DialogActions>
      </form>
    </Dialog>
  );
}

export default EnableMfaModal;
