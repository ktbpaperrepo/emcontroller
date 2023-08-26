import random
import string
import time

DATA_FILE_NAME_PREFIX = "data/resp_data_"
DATA_FILE_NAME_SUFFIX = ".csv"
DATA_FILE_RANDOM_LENGTH = 10


# use lowercase letters and digits to randomly generate a data file name
def gen_data_file_name() -> str:
    letters_digits = string.ascii_lowercase + string.digits
    data_file_name = DATA_FILE_NAME_PREFIX  # prefix
    data_file_name += str(time.time()) + "_"  # timestamp

    # random part
    for i in range(DATA_FILE_RANDOM_LENGTH):
        data_file_name += random.choice(letters_digits)

    data_file_name += DATA_FILE_NAME_SUFFIX  # suffix
    return data_file_name
