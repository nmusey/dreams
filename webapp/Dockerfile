FROM node:14
WORKDIR /usr/app

COPY package* .
RUN npm install

COPY . .

EXPOSE 3000
CMD ["npm", "run", "dev"]
