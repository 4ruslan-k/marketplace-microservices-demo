import * as Joi from 'joi';

const configSchema = Joi.object({
  port: Joi.number().required(),
  natsUri: Joi.string().required(),
  database: {
    host: Joi.string().required(),
    port: Joi.number().required(),
    username: Joi.string().required(),
    password: Joi.string().required(),
    name: Joi.string().required(),
    debugMode: Joi.boolean().default(false),
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
      debugMode: process.env.DATABASE_DEBUG_MODE === 'true' || false,
    },
    natsUri: process.env.NATS_URI,
  };

  const validate = () => {
    const { error } = configSchema.validate(config, { abortEarly: false });
    if (error) throw error;
  };

  validate();

  return config;
};
