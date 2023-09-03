import * as Joi from 'joi';

const configSchema = Joi.object({
  port: Joi.number().required(),
  database: {
    host: Joi.string().required(),
    port: Joi.number().required(),
    username: Joi.string().required(),
    password: Joi.string().required(),
    name: Joi.string().required(),
  },
});

export default () => {
  const config = {
    port: parseInt(process.env.PORT, 10) || 3000,
    database: {
      host: process.env.DATABASE_HOST,
      port: parseInt(process.env.DATABASE_PORT, 10),
      username: process.env.DATABASE_USERNAME,
      password: process.env.DATABASE_PASSWORD,
      name: process.env.DATABASE_NAME,
    },
  };

  const validate = () => {
    const { error } = configSchema.validate(config, { abortEarly: false });
    if (error) throw error;
  };

  validate();

  return config;
};
